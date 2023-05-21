package ro.cryptogen.cryptogenexpose.service

import org.springframework.stereotype.Service
import ro.cryptogen.cryptogenexpose.model.CryptoMessage
import ro.cryptogen.cryptogenexpose.repository.CryptoRepository

@Service
class CryptoHistoryService(private val cryptoRepository: CryptoRepository) {

    fun getCryptoHistory(): MutableList<CryptoMessage> {
        return cryptoRepository.findAll()
    }
}