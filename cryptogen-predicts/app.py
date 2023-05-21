# Importing the necessary packages
from flask import Flask, request
import numpy as np
import pandas as pd
import datetime as dt
import yfinance as yf
import tensorflow as tf
from sklearn.preprocessing import MinMaxScaler
import aio_pika
import asyncio
import json

# Initialize the Flask application
app = Flask(__name__)

# Global variables
prediction_days = 7
future_day = 3
against_currency = 'USD'
scaler = MinMaxScaler(feature_range=(0, 1))

# Send message to RabbitMQ
async def send_to_rabbitmq(message):
    connection = await aio_pika.connect_robust("amqp://guest:guest@localhost/")

    async with connection:
        channel = await connection.channel()

        exchange = await channel.declare_exchange('crypto_exchange', aio_pika.ExchangeType.FANOUT)

        message_json = json.dumps(message)
        message = aio_pika.Message(body=message_json.encode())

        await exchange.publish(message, routing_key='crypto_queue')

        print(" [x] Sent %r" % message.body)

def train_model(symbol):
    model_file = f'{symbol}_model.h5'

    start = dt.datetime(2016, 1, 1)
    end = dt.datetime.now()

    data = yf.download(f'{symbol}-{against_currency}', start=start, end=end)

    scaled_data = scaler.fit_transform(data['Close'].values.reshape(-1, 1))

    x_train, y_train = [], []

    for x in range(prediction_days, len(scaled_data) - future_day):
        x_train.append(scaled_data[x - prediction_days:x, 0])
        y_train.append(scaled_data[x + future_day, 0])

    x_train, y_train = np.array(x_train), np.array(y_train)
    x_train = np.reshape(x_train, (x_train.shape[0], x_train.shape[1], 1))

    model = tf.keras.Sequential()

    model.add(tf.keras.layers.LSTM(units=50, return_sequences=True, input_shape=(x_train.shape[1], 1)))
    model.add(tf.keras.layers.Dropout(0.2))
    model.add(tf.keras.layers.LSTM(units=50, return_sequences=True))
    model.add(tf.keras.layers.Dropout(0.2))
    model.add(tf.keras.layers.LSTM(units=50))
    model.add(tf.keras.layers.Dropout(0.2))
    model.add(tf.keras.layers.Dense(units=1))

    model.compile(optimizer='adam', loss='mean_squared_error')
    model.fit(x_train, y_train, epochs=25, batch_size=32)

    model.save(model_file) # Save the trained model

def predict_price(symbol):
    model_file = f'{symbol}_model.h5'

    test_start = dt.datetime(2022, 1, 1)
    test_end = dt.datetime.now()

    test_data = yf.download(f'{symbol}-{against_currency}', start=test_start, end=test_end)
    actual_prices = test_data['Close'].values

    # Load the old data
    start = dt.datetime(2016, 1, 1)
    end = dt.datetime.now()
    old_data = yf.download(f'{symbol}-{against_currency}', start=start, end=end)

    total_dataset = pd.concat((old_data['Close'], test_data['Close']), axis=0)

    model_inputs = total_dataset[len(total_dataset) - len(test_data) - prediction_days:].values
    model_inputs = model_inputs.reshape(-1, 1)
    model_inputs = scaler.fit_transform(model_inputs)

    x_test = []

    for x in range(prediction_days, len(model_inputs)):
        x_test.append(model_inputs[x - prediction_days:x, 0])

    x_test = np.array(x_test)
    x_test = np.reshape(x_test, (x_test.shape[0], x_test.shape[1], 1))

    # Load the trained model
    model = tf.keras.models.load_model(model_file)

    prediction_prices = model.predict(x_test)
    prediction_prices = scaler.inverse_transform(prediction_prices)

    # Convert the prediction prices to list
    predictions_list = prediction_prices[-3:].flatten().tolist()
    actual_prices_list = actual_prices[-3:].tolist()

    # Send predictions and actual prices to RabbitMQ
    asyncio.run(send_to_rabbitmq({
        'predictions': predictions_list,
        'actual_prices': actual_prices_list
    }))

    return predictions_list, actual_prices_list

# Define the route for the POST request
@app.route('/crypto', methods=['POST'])
def process_request():
    request_data = request.get_json()
    symbol = request_data['symbol']
    train = request_data.get('train', False)

    if train:
        train_model(symbol)

    predictions, actual_prices = predict_price(symbol)

    return {
        'message': f'Prediction for {symbol} processed successfully!',
        'predictions': predictions,
        'actual_prices': actual_prices
    }

if __name__ == "__main__":
    app.run(port=5100)