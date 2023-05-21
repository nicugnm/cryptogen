package ro.cryptogen.cryptogenexpose.controller

import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController
import ro.cryptogen.cryptogenexpose.model.CryptoMessage
import ro.cryptogen.cryptogenexpose.service.CryptoHistoryService
import ro.cryptogen.cryptogenexpose.service.MessageProcessingService

@RestController
@RequestMapping("/api/crypto")
class MessageController(private val messageProcessingService: MessageProcessingService,
    private val cryptoHistoryService: CryptoHistoryService) {

    @GetMapping("/last-message")
    fun getLastMessage(): String {
        return messageProcessingService.lastMessage
    }

    @GetMapping("/history")
    fun getHistoryMessages() : MutableList<CryptoMessage> {
        return cryptoHistoryService.getCryptoHistory()
    }
}