package ro.cryptogen.cryptogenexpose.payload

import lombok.AllArgsConstructor
import lombok.Builder

@Builder
@AllArgsConstructor
data class SymbolRequest(var symbol: String, var train: Boolean) {
}