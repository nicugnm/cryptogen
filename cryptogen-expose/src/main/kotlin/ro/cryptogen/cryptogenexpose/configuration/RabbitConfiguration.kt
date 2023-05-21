import org.springframework.amqp.core.Binding
import org.springframework.amqp.core.BindingBuilder
import org.springframework.amqp.core.Exchange
import org.springframework.amqp.core.ExchangeBuilder
import org.springframework.amqp.core.Queue
import org.springframework.amqp.core.QueueBuilder
import org.springframework.amqp.rabbit.config.SimpleRabbitListenerContainerFactory
import org.springframework.amqp.rabbit.connection.CachingConnectionFactory
import org.springframework.amqp.rabbit.listener.adapter.MessageListenerAdapter
import org.springframework.beans.factory.annotation.Value
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import ro.cryptogen.cryptogenexpose.service.MessageListenerService

@Configuration
class RabbitMqConfiguration {

    @Value("\${spring.rabbitmq.host}")
    private lateinit var host: String

    @Value("\${spring.rabbitmq.port}")
    private var port: Int = 0

    @Value("\${spring.rabbitmq.username}")
    private lateinit var username: String

    @Value("\${spring.rabbitmq.password}")
    private lateinit var password: String

    @Bean
    fun connectionFactory(): CachingConnectionFactory {
        var connectionFactory = CachingConnectionFactory(host, port)
        connectionFactory.username = username
        connectionFactory.setPassword(password)
        return connectionFactory
    }

    @Bean
    fun listenerAdapter(messageListenerService: MessageListenerService): MessageListenerAdapter {
        return MessageListenerAdapter(messageListenerService, "receiveMessage")
    }

    @Bean
    fun rabbitListenerContainerFactory(): SimpleRabbitListenerContainerFactory {
        val factory = SimpleRabbitListenerContainerFactory()
        factory.setConnectionFactory(connectionFactory())
        return factory
    }

    @Bean
    fun exchange(): Exchange {
        return ExchangeBuilder.fanoutExchange("crypto_exchange").durable(true).build()
    }

    @Bean
    fun queue(): Queue {
        return QueueBuilder.nonDurable().build()
    }

    @Bean
    fun binding(queue: Queue, exchange: Exchange): Binding {
        return BindingBuilder.bind(queue).to(exchange).with("*").noargs()
    }
}