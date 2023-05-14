package ro.cryptogen.cryptogenexpose.configuration

import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import java.net.http.HttpClient

@Configuration
class InjectionConfig {

    @Bean
    fun httpClient(): HttpClient {
        return HttpClient.newHttpClient()
    }
}