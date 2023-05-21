package ro.cryptogen.cryptogenexpose.service

import org.springframework.amqp.rabbit.annotation.Queue
import org.springframework.amqp.rabbit.annotation.RabbitListener
import org.springframework.context.ApplicationEventPublisher
import org.springframework.stereotype.Component
import ro.cryptogen.cryptogenexpose.events.MessageReceivedEvent
import ro.cryptogen.cryptogenexpose.model.CryptoMessage
import ro.cryptogen.cryptogenexpose.repository.CryptoRepository

@Component
class MessageListenerService(private val eventPublisher: ApplicationEventPublisher,
    private val cryptoRepository: CryptoRepository) {

    @RabbitListener(queuesToDeclare = [Queue("crypto_exchange")])
    fun receiveMessage(message: String) {
        eventPublisher.publishEvent(MessageReceivedEvent(this, message))
        cryptoRepository.save(CryptoMessage(content = message))
    }
}