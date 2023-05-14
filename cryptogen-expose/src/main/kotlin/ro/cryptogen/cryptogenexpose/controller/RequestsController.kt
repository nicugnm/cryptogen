package ro.cryptogen.cryptogenexpose.controller

import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController
import ro.cryptogen.cryptogenexpose.service.CryptoTriggerService

@RestController
@RequestMapping("/api/triggers")
class RequestsController(private val cryptoTriggerService: CryptoTriggerService) {

    @PostMapping("/{symbol}")
    fun requestPrediction(@PathVariable(name = "symbol") symbol: String) {
        cryptoTriggerService.createPredictionRequest(symbol)
    }
}