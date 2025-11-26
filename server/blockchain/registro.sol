pragma solidity ^0.8.0;

contract Registro {
    struct Card {
        string element;
        string cardType;
        string id;
        address dono;
    }

    mapping(string => Card) public cards;

    function registrarCardParaJogador(
        address _jogadorDestino,
        string memory _id, 
        string memory _element, 
        string memory _cardType
    ) public {
        require(cards[_id].dono == address(0), "ID ja em uso");
        cards[_id] = Card(_element, _cardType, _id, _jogadorDestino);
    }

    function transferirCard(string memory _id, address _novoDono) public {
        require(cards[_id].dono == msg.sender, "Voce nao e o dono");
        cards[_id].dono = _novoDono;
    }
}