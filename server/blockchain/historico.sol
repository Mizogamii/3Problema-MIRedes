// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Historico {
    // 0 = EMPATE
    // 1 = VITORIA_JOGADOR_1
    // 2 = VITORIA_JOGADOR_2
    //enum Resultado { EMPATE, VITORIA_J1, VITORIA_J2 }

    //event NovaPartida(string jogador1, string jogador2, Resultado resultado, uint256 data);

    event NovaPartida(string jogador1, string jogador2, string resultado, uint256 data);


    function registrarPartida(string memory _j1, string memory _j2, string _resultado) public {
        //require(_resultado <= 2, "Resultado invalido (use 0, 1 ou 2)");
        
        //emit NovaPartida(_j1, _j2, Resultado(_resultado), block.timestamp);
        emit NovaPartida(_j1, _j2, _resultado, block.timestamp);
    }
}