pragma solidity ^0.8.0;

contract HistoricoDePartidas {
    
    event NovaPartida(string vencedor, string perdedor, string placar, uint256 data);

    function registrarResultado(string memory _vencedor, string memory _perdedor, string memory _placar) public {
        // 2. Em vez de salvar numa variável, nós "emitimos" o evento
        // Isso grava no log do bloco permanentemente.
        emit NovaPartida(_vencedor, _perdedor, _placar, block.timestamp);
    }
}