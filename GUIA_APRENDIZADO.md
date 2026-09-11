# 📘 Relatório Técnico & Guia de Aprendizado: Construção do Paidégua Game em Go

**Data:** 10 de Setembro de 2026  
**Criador & Desenvolvedor:** Lucivaldo Junior (Luci Junior)  
**Co-criação & Arquitetura Técnica:** Nexus AI Ecosystem (sob liderança de Lucy — Tech Lead Sênior)  
**Projeto:** `paidegua-game` (`github.com/luci-jr/paidegua-game`)  
**Foco:** Engenharia de Software em Go, Padrão de Projeto Real e Preparação para Entrevista Backend Júnior  

---

## 🎯 Objetivo Deste Documento
Registrar passo a passo toda a jornada de engenharia percorrida na criação do jogo **Paidégua Game**, desde o primeiro comando no terminal e a resolução dos erros iniciais de compilação até a arquitetura modular corporativa final com áudio procedural e progressão de fases em Belém do Pará.

---

## 🧭 Índice do Aprendizado
1. [Fase 0: A Base do Go e o Módulo (`go.mod`)](#fase-0-a-base-do-go-e-o-módulo-gomod)
2. [Fase 1: O Primeiro Modelo Físico & Structs Idiomáticas](#fase-1-o-primeiro-modelo-físico--structs-idiomáticas)
3. [Fase 2: A Transição Gráfica com Ebitengine 2D](#fase-2-a-transição-gráfica-com-ebitengine-2d)
4. [Fase 3: Controles de Teclado & O Dilema dos Obstáculos](#fase-3-controles-de-teclado--o-dilema-dos-obstáculos)
5. [Fase 4: Vidas, i-Frames e Respiro Visual (HUD)](#fase-4-vidas-i-frames-e-respiro-visual-hud)
6. [Fase 5: A Grande Refatoração Arquitetural (`cmd/` e `internal/`)](#fase-5-a-grande-refatoração-arquitetural-cmd-e-internal)
7. [Fase 6: O Salto Anatômico da Onça e o Game Feel](#fase-6-o-salto-anatômico-da-onça-e-o-game-feel)
8. [Fase 7: Identidade Cultural de Belém & Síntese Procedural de Áudio](#fase-7-identidade-cultural-de-belém--síntese-procedural-de-áudio)
9. [Fase 8: Mercado de São Brás, Mobile Touch & Ponte WebAssembly](#fase-8-mercado-de-são-brás-mobile-touch--ponte-webassembly)
10. [Guia de Defesa Técnica para a Entrevista](#guia-de-defesa-técnica-para-a-entrevista)

---

## 🧱 Fase 0: A Base do Go e o Módulo (`go.mod`)

### O que aprendemos:
* **O Erro `expected 'package', found 'EOF'`:**
  * Ocorre quando um arquivo `.go` está com 0 bytes ou não foi salvo no editor. No Go, todo arquivo deve começar explicitamente com `package <nome>`.
* **A Necessidade do `go.mod`:**
  * No Go moderno (1.16+), o comando `go get` não funciona solto. O Go exige um módulo inicializado para rastrear dependências com controle estrito de versão:
    ```bash
    go mod init github.com/luci-jr/paidegua-game
    go get github.com/hajimehoshi/ebiten/v2
    ```
  * O arquivo `go.mod` equivale ao `pom.xml` do Maven (Java) ou ao `package.json` (Node.js). O `go.sum` garante a integridade criptográfica dos pacotes baixados.

---

## 🐾 Fase 1: O Primeiro Modelo Físico & Structs Idiomáticas

Começamos simulando a física em texto puro no terminal.

### Conceitos-chave de Go:
1. **Separação de Dados e Comportamento:**
   * Go não tem classes (`class`). Usamos `type Struct` para os dados.
2. **Receivers de Ponteiro (`*`):**
   * Em `func (o *Onca) Jump()`, o `*Onca` indica que a função opera na referência da memória real da instância, e não em uma cópia.
3. **Física da Gravidade:**
   * O salto aplica um impulso negativo (`VelocityY = -7.5`), e a cada frame somamos a gravidade (`VelocityY += 0.4`), criando a trajetória parabólica real.

---

## 🖼️ Fase 2: A Transição Gráfica com Ebitengine 2D

Para sair do terminal e abrir uma janela de verdade a 60 FPS, escolhemos a biblioteca padrão da indústria em Go: **Ebitengine (v2)**.

### A Interface Implícita `ebiten.Game`:
No Java você declararia `class Jogo implements Game`. No Go, **interfaces são satisfeitas implicitamente**. Basta a sua struct ter estes 3 métodos:
1. **`Update() error`**: Executado 60 vezes por segundo (TPS fixo) para calcular física, colisões e ler teclado.
2. **`Draw(screen *ebiten.Image)`**: Chamado a cada frame para renderizar sprites e pixels na tela.
3. **`Layout(outsideW, outsideH int) (screenW, screenH int)`**: Define a resolução lógica retrô fixa (ex: `340x210`), garantindo que a tela escale perfeitamente em monitores Full HD ou 4K sem distorcer pixels.

---

## 🎮 Fase 3: Controles de Teclado & O Dilema dos Obstáculos

### O problema resolvido:
Se o jogador só precisa pular, a tecla **Seta para Baixo** (abaixar) perde a utilidade.

### A solução de Game Design:
* Criamos dois tipos de obstáculos no loop:
  1. **Obstáculo no Chão (Paneiro de Açaí):** Ocupa a base do solo. Se abaixar, bate nele. **Obrigatório PULAR!**
  2. **Obstáculo no Ar (Urubu voando a meia altura):** Voa na altura da cabeça da onça. Se pular ou ficar em pé, colide. **Obrigatório ABAIXAR!**
* **Bounding Boxes Dinâmicas:** Ao apertar a Seta para Baixo, o hitbox da onça diminui de 22px para 12px de altura, permitindo que as aves passem por cima das suas costas.

---

## 💖 Fase 4: Vidas, i-Frames e Respiro Visual (HUD)

### Desafios de Engenharia:
1. **O problema do dano contínuo:** Sem tratamento, colidir com um obstáculo causaria dano em todos os 60 frames por segundo, matando a onça instantaneamente.
   * **Solução (Invincibility Frames / i-frames):** Ao tomar dano, o jogo ativa `invincibleTicks = 60` (1 segundo). Durante esse tempo, a onça pisca e fica imune a novos danos.
2. **Corações em Pixel Art:**
   * Criamos a função matemática `drawHeart(screen, x, y, filled)` que desenha corações cheios vermelhos brilhantes e, após sofrer dano, corações vazados (apenas contorno).
3. **Pausa com a tecla `Esc`:** Congela o loop físico e exibe a tela de pausa com overlay translúcido.

---

## 🏛️ Fase 5: A Grande Refatoração Arquitetural (`cmd/` e `internal/`)

Ao atingir 350 linhas de código no `main.go`, aplicamos o **Standard Go Project Layout**, transformando o monólito em pacotes coesos:

```text
paidegua-game/
├── cmd/runner/main.go        # Ponto de entrada fino (Thin Main)
├── internal/
│   ├── audio/audio.go        # Gerenciamento de som
│   ├── entities/onca.go      # Entidade Onça (física, sprites, partículas)
│   ├── entities/obstacle.go  # Entidade Obstáculos (chão e ar)
│   ├── scenery/ver_o_peso.go # Cenários em camadas
│   ├── ui/hud.go             # Interface de usuário e textos
│   └── game/game.go          # Orquestrador central (Engine)
```

### Por que essa separação é crucial?
* **Regra do `internal/` no Go:** Pacotes dentro de `internal/` são protegidos pelo compilador do Go. Nenhum projeto externo consegue importá-los, garantindo privacidade arquitetural.
* **Manutenibilidade:** Cada arquivo tem responsabilidade única (SRP - Single Responsibility Principle).

---

## 🐆 Fase 6: O Salto Anatômico da Onça e o Game Feel

Substituímos os blocos geométricos simplórios por um desenho com anatomia real da **Onça-Pintada Brasileira**:
* **Rosetas Autênticas:** Manchas escuras com o centro ocre-alaranjado.
* **Ventre e Queixo Claros:** Destacando a silhueta muscular do felino.
* **Olhos Amendoados Verdes:** Com pupila vertical e brilho.
* **Ciclo de Galope em 4 Quadros:**
  * Fase 0: Alongamento máximo no ar.
  * Fase 1: Impacto das patas dianteiras.
  * Fase 2: Recolhimento em mola sob o ventre.
  * Fase 3: Impulso traseiro poderoso.
* **Game Feel Moderno:**
  * **Pulo com Altura Variável:** Soltar o botão corta a velocidade vertical para saltos curtos.
  * **Coyote Time:** Tolerância de alguns milissegundos para saltar após sair do chão.
  * **Sistema de Partículas:** Poeira procedural sob as patas na aterrissagem e no rastejo.
  * **Screen Shake:** Tremor suave da tela ao sofrer impacto.
  * **Hit-stop & Balão "Égua mano!...":** Ao perder um coração, a simulação congela brevemente por 22 ticks (~0.35s), toca um efeito sonoro cômico característico e exibe um balão de fala retrô com a gíria paraense sobre a onça.

---

## 🌴 Fase 7: Identidade Cultural de Belém & Síntese Procedural de Áudio

### 1. As 3 Fases Históricas:
* **Fase 1 (0 a 100 pts): Mercado do Ver-o-Peso**
  * Mercado de Ferro com cúpulas azuis, Baía do Guajará, Paneiro de Açaí e Urubus.
* **Fase 2 (101 a 200 pts): Estação das Docas**
  * Armazéns ingleses vermelhos, guindaste amarelo portuário, Barril das Docas e Gaivotas.
* **Fase 3 (201 a 300 pts): Theatro da Paz & Mangueiras**
  * Fachada neoclássica nobre, mangueiras frondosas com mangas maduras e calçada portuguesa.

### 2. Menu de Pausa & Mute Dinâmico:
* **Pausa Interativa com `Esc`:** Interrompe a física e a música, abrindo modal retrô com cursor dourado.
* **Navegação com `Setas Cima/Baixo` e `Enter`:** Permite continuar, reiniciar a aventura desde a Fase 1, alternar o som entre `LIGADO` e `MUDO` (com indicador visual no HUD), abrir os créditos e fechar o jogo de forma limpa via `ebiten.Termination`.

### 3. Tela de Início (Title Screen) & Créditos Finais:
* **Início Profissional:** Ao abrir o jogo, é apresentada a Tela Inicial com os metadados técnicos: **Linguagem Utilizada: Go (Golang 1.22+)**, **Desenvolvedor: Lucivaldo Junior (Luci)**, **Game Engine: Ebitengine (v2)** e instruções de comando.
* **Créditos Finais:** Exibidos automaticamente na conclusão da Fase 3 (vitória da expedição) e acessíveis a qualquer instante com a tecla `C` ou pelo menu de pausa, destacando a autoria e stack técnica do projeto.

### 4. Síntese de Áudio Procedural (Zero Arquivos Externos):
* **Trilha de Carimbó 8-bit:** Batida sincopada de curimbó grave (82Hz), maracas rítmicas com ruído filtrado no contratempo e melodia contagiante tocando em loop contínuo via `audio.NewInfiniteLoop`.
* **Rugido Felino da Onça (*Feline Roar*):** Ao saltar, dispara uma frequência gutural (95Hz a 125Hz) modulada por tremolo de garganta (42Hz) e ruído de respiração felina.

---

## 🏛️ Fase 8: Mercado de São Brás, Mobile Touch & Ponte WebAssembly

### 1. Passagem Histórica pelo Mercado de São Brás Atual:
* **Homenagem ao Patrimônio de Belém:** O jogo inicia com uma passagem cinematográfica pelo **Mercado de São Brás revitalizado (1904 - 2024)**, exibindo a imponente fachada histórica restaurada, a torre do relógio, praça arborizada e quiosques tradicionais de açaí e tapioca.
* **Experiência do Jogador:** Permanece visível por ~15 segundos para contemplação, ou avança imediatamente ao receber qualquer clique de mouse, toque na tela do celular ou tecla pressionada.

### 2. Menu Inicial Compacto & Visibilidade do Cenário:
* **Dimensões Enxutas (`195 x 108 px`):** O menu foi reduzido para não cobrir a tela, utilizando fundo translúcido suave.
* **Cenário Dinâmico:** Permite ver todo o visual histórico e a onça correndo em galope animado ao fundo enquanto o menu está aberto.

### 3. Movimentação Horizontal Bidirecional da Onça:
* A onça agora pode **adiantar** (`→` / `D`) e **recuar** (`←` / `A`) na tela horizontalmente.
* **Contenção Matemática (*Clamping*):** Limites dinâmicos entre `X = 15.0` e `X = 230.0` para manter o posicionamento seguro em relação aos obstáculos.

### 4. Arquitetura Mobile & A Ponte de Memória JavaScript-Go (`syscall/js`):
* **O Desafio do Mobile WebAssembly:** Navegadores em smartphones bloqueiam eventos de teclado sintéticos (`isTrusted: false`) enviados via JavaScript.
* **A Solução de Engenharia:**
  * Criamos uma ponte direta via `syscall/js`: o JavaScript altera o estado do objeto global `window._virtualKeys`, e o Go lê esse estado diretamente da memória a cada tick via o módulo [`internal/game/input_js.go`](file:///home/lucivaldo-junior/Documentos/GitHub/projetos/paidegua-game/internal/game/input_js.go) com *Build Tags* (`//go:build js && wasm`).
  * Para compilações desktop nativas, criamos o fallback elegante [`internal/game/input_other.go`](file:///home/lucivaldo-junior/Documentos/GitHub/projetos/paidegua-game/internal/game/input_other.go) (`//go:build !(js && wasm)`).
  * Lemos também toques físicos nativos no canvas com `ebiten.AppendTouchIDs` e `ebiten.TouchPosition`.

### 5. Balão de Morte Autêntico de Belém:
* **Dano Intermediário:** A onça exibe o balão cômico regional com som: `"EGUA MANO!..."`.
* **Morte Final (Game Over):** A onça exibe o clássico balão paraense: `"Levei o farelo mano, mancada!"` apontando para a sua cabeça, com a janela de pontuação reposicionada no topo da tela.

---

## 🎤 Guia de Defesa Técnica para a Entrevista (Jungle Gaming)

Ao apresentar este projeto em sua entrevista para **Backend Go Júnior**, destaque estes pontos:

1. **"Por que você escolheu Go para este projeto?"**
   * *"Go entrega altíssima performance com baixo consumo de memória, compilação estática nativa sem necessidade de JVM ou interpretador, e tratamento rigoroso e determinístico de concorrência e loops temporais."*
2. **"Como você organizou a arquitetura do código?"**
   * *"Segui o Standard Go Project Layout. Isolei as entidades de negócio, física e áudio dentro da pasta `internal/` para garantir encapsulamento nativo pelo compilador, mantendo o `cmd/runner/main.go` enxuto com apenas a inicialização de dependências."*
3. **"Como funcionam as colisões e a física?"**
   * *"Implementei um loop de física desacoplado no `Update()` a 60 ticks por segundo, utilizando detecção de colisão AABB (Axis-Aligned Bounding Box) com caixas dinâmicas que variam conforme o estado da onça (em pé, pulando, agachada ou movendo-se lateralmente) e tolerância temporal com Coyote Time."*
4. **"Como você resolveu a compatibilidade entre Desktop e Web/Mobile?"**
   * *"Utilizei Build Tags (`//go:build js && wasm` vs `//go:build !(js && wasm)`) para separar a camada de input. No navegador e celular, o Go se comunica diretamente com o JavaScript via o pacote `syscall/js`, permitindo que botões virtuais na tela alterem o estado de memória diretamente sem depender de eventos sintéticos de teclado bloqueados por navegadores móveis."*
5. **"Como foi resolvido o áudio sem travar a thread principal?"**
   * *"Utilizei streams PCM estéreo a 44.1kHz sintetizados matematicamente na inicialização e tocados assincronamente via o subsistema de áudio do Ebitengine, com loop infinito sem fim em buffer de memória."*
6. **"Como foi o processo de desenvolvimento e o uso de IA?"**
   * *"O projeto foi concebido e desenvolvido por mim (Lucivaldo Junior) em co-criação com o **Nexus AI Ecosystem**, um squad autônomo de múltiplos agentes de IA que eu mesmo desenvolvi e configurei. A IA atuou como pair programmer sênior (sob a liderança de Lucy - Tech Lead & Arquiteta), auxiliando na governança arquitetural, benchmarking de física, compilação WebAssembly e síntese de áudio procedural."*

---

> **Status Final:** Projeto 100% funcional, modularizado, documentado, compilado para WebAssembly e Desktop Nativo, co-criado com Nexus e aprovado com distinção para apresentação técnica.
