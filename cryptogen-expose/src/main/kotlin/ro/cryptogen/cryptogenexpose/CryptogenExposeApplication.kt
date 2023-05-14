package ro.cryptogen.cryptogenexpose

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.boot.runApplication
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RestController
import ro.cryptogen.cryptogenexpose.configuration.ApiProperties

@SpringBootApplication
@EnableConfigurationProperties(ApiProperties::class)
class CryptogenExposeApplication

fun main(args: Array<String>) {
	runApplication<CryptogenExposeApplication>(*args)
}