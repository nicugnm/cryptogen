package ro.cryptogen.cryptogenexpose.controller

import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController
import ro.cryptogen.cryptogenexpose.payload.SymbolRequest
import ro.cryptogen.cryptogenexpose.service.CryptoTriggerService

@RestController
@RequestMapping("/api/triggers")
class RequestsController(private val cryptoTriggerService: CryptoTriggerService) {

    @PostMapping
    fun requestPrediction(@RequestBody symbolRequest: SymbolRequest) {
        cryptoTriggerService.createPredictionRequest(symbolRequest)
    }
}