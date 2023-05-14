package ro.cryptogen.cryptogenexpose.service

import org.springframework.stereotype.Service
import ro.cryptogen.cryptogenexpose.configuration.ApiProperties
import java.net.URI
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpRequest.BodyPublishers
import java.net.http.HttpResponse

@Service
class CryptoTriggerService(
    private val httpClient: HttpClient,
    private val apiProperties: ApiProperties
) {
    fun createPredictionRequest(symbol: String) {
        val httpRequest: HttpRequest = HttpRequest.newBuilder()
            .uri(URI.create(apiProperties.url))
            .header("Content-Type", "application/json")
            .POST(BodyPublishers.ofString(symbol))
            .build()

        httpClient.send(httpRequest, HttpResponse.BodyHandlers.ofString())
    }
}




