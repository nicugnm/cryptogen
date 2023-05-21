package ro.cryptogen.cryptogenexpose.service

import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper
import org.springframework.stereotype.Service
import ro.cryptogen.cryptogenexpose.configuration.ApiProperties
import ro.cryptogen.cryptogenexpose.payload.SymbolRequest
import java.net.URI
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpResponse

@Service
class CryptoTriggerService(
    private val httpClient: HttpClient,
    private val apiProperties: ApiProperties
) {
    fun createPredictionRequest(symbolRequest: SymbolRequest) {
        val mapper = jacksonObjectMapper()
        val httpRequest: HttpRequest = HttpRequest.newBuilder()
            .uri(URI.create(apiProperties.url))
            .header("Content-Type", "application/json")
            .POST(HttpRequest.BodyPublishers.ofString(mapper.writeValueAsString(symbolRequest)))
            .build()

        httpClient.send(httpRequest, HttpResponse.BodyHandlers.ofString())
    }
}
