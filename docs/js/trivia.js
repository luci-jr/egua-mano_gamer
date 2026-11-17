/**
 * ============================================================================
 * Égua Mano! Gamer - Curiosidades Históricas e Culturais de Belém do Pará
 * Clean Architecture: Módulo de Conteúdo e Ambientação Regional
 * ============================================================================
 */

(function () {
  'use strict';

  const belemTrivia = [
    "O Mercado do Ver-o-Peso (1627) é a maior feira a céu aberto da América Latina!",
    "Belém é conhecida carinhosamente como a 'Cidade das Mangueiras', plantadas no século XIX.",
    "O açaí paraense tradicional é consumido puro na cuia com peixe frito e farinha d'água!",
    "O Theatro da Paz (1878) é uma joia neoclássica inspirada no Teatro Scala de Milão.",
    "O Mercado de São Brás (1911) foi revitalizado como grande polo cultural e gastronômico!",
    "A Estação das Docas preserva armazéns de ferro ingleses de 1897 e guindastes portuários históricos.",
    "O Círio de Nazaré reúne mais de 2 milhões de fiéis nas ruas no 2º domingo de outubro!",
    "O Tacacá leva tucupi e folhas de jambu, que amortecem suavemente a boca graças ao espilantol.",
    "A expressão paraense 'Paidégua' é usada para descrever algo excelente, sensacional e autêntico!",
    "Em Belém, é tradicional marcar encontros 'antes ou depois da chuva da tarde'.",
    "A Ilha do Combú fica em frente a Belém e produz chocolates artesanais com cacau nativo.",
    "A Onça-Pintada é o maior felino das Américas e rainha suprema das florestas amazônicas!",
    "O Carimbó, ritmo tradicional paraense marcado pelo curimbó, é Patrimônio Cultural do Brasil!",
    "O Jacaré-Açu pode ultrapassar 4,5 metros de comprimento nos igarapés e rios do Pará!",
    "Na infância do Pará e da Amazônia, o estilingue é chamado carinhosamente de 'Baladeira'!"
  ];

  let triviaIdx = 0;
  const triviaElem = document.getElementById('trivia-text');

  if (triviaElem) {
    setInterval(() => {
      triviaElem.style.opacity = '0';
      setTimeout(() => {
        triviaIdx = (triviaIdx + 1) % belemTrivia.length;
        triviaElem.textContent = belemTrivia[triviaIdx];
        triviaElem.style.opacity = '1';
      }, 300);
    }, 6000);
  }

  window.BelemTrivia = {
    getAll: () => belemTrivia,
    getCurrent: () => belemTrivia[triviaIdx]
  };
})();
