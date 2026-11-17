/**
 * ============================================================================
 * Égua Mano! Gamer - Carregador Inteligente de WebAssembly
 * Clean Architecture: Adaptador de Inicialização de Runtime e Binários
 * ============================================================================
 * Carrega o binário compactado em gzip (game.dat) via DecompressionStream,
 * realizando fallback automático para o arquivo descompactado (game.wasm)
 * caso necessário.
 */

(function () {
  'use strict';

  async function loadCompressedWasm() {
    const loadingElem = document.getElementById('loading');
    const go = new Go();

    let response = await fetch('game.dat');
    let isGzip = true;

    // Se game.dat não existir localmente, tenta carregar o game.wasm direto
    if (!response.ok) {
      response = await fetch('game.wasm');
      isGzip = false;
      if (!response.ok) {
        throw new Error('HTTP 404: nem game.dat nem game.wasm foram encontrados');
      }
    }

    let wasmBytes;
    if (isGzip) {
      if (typeof DecompressionStream !== 'undefined') {
        // Navegadores modernos: descompressão nativa por stream
        const ds = new DecompressionStream('gzip');
        const decompressed = response.body.pipeThrough(ds);
        wasmBytes = await new Response(decompressed).arrayBuffer();
      } else {
        // Fallback: descompressão via cabeçalho HTTP simulado
        const compressed = await response.arrayBuffer();
        const blob = new Blob([compressed]);
        const dsResponse = new Response(blob, {
          headers: { 'Content-Encoding': 'gzip', 'Content-Type': 'application/wasm' }
        });
        wasmBytes = await dsResponse.arrayBuffer();
      }
    } else {
      wasmBytes = await response.arrayBuffer();
    }

    const result = await WebAssembly.instantiate(wasmBytes, go.importObject);

    if (loadingElem) {
      loadingElem.style.display = 'none';
    }

    go.run(result.instance);
    window.focus();
  }

  window.addEventListener('load', () => {
    loadCompressedWasm().catch((err) => {
      console.error('[WasmLoader] Erro ao carregar o jogo:', err);
      const loadingElem = document.getElementById('loading');
      if (loadingElem) {
        loadingElem.innerHTML = `
          <div style="background: rgba(15, 23, 42, 0.9); padding: 18px; border-radius: 12px; border: 1px solid #ef4444; max-width: 80vw;">
            <p style="color: #ef4444; font-weight: bold; margin-bottom: 8px;">Erro ao carregar o jogo</p>
            <p style="color: #fca5a5; font-size: 13px;">${err.message || err}</p>
          </div>
        `;
      }
    });
  });

  window.WasmLoader = {
    load: loadCompressedWasm
  };
})();
