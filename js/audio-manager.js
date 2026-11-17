/**
 * ============================================================================
 * Égua Mano! Gamer - Gerenciador de Áudio & Ciclo de Vida da Aba/App
 * Clean Architecture: Módulo de Infraestrutura de Som e Foco
 * ============================================================================
 * Silencia instantaneamente o jogo quando o usuário sai da aba ou minimiza
 * (ex: responder mensagem no WhatsApp, alternar janelas ou bloquear a tela).
 */

(function () {
  'use strict';

  const activeAudioContexts = [];
  const OrigAudioContext = window.AudioContext || window.webkitAudioContext;

  if (OrigAudioContext) {
    try {
      const AudioCtxProxy = new Proxy(OrigAudioContext, {
        construct(target, args) {
          const ctx = new target(...args);
          activeAudioContexts.push(ctx);
          return ctx;
        }
      });
      window.AudioContext = AudioCtxProxy;
      if (window.webkitAudioContext) {
        window.webkitAudioContext = AudioCtxProxy;
      }
    } catch (e) {
      console.warn('[AudioManager] Falha ao interceptar AudioContext:', e);
    }
  }

  let isAppInForeground = true;

  function suspendAppAudio() {
    if (!isAppInForeground) return;
    isAppInForeground = false;

    if (window._virtualKeys) {
      window._virtualKeys.PageHidden = true;
      window._virtualKeys.JustUnfocused = true;
    }

    activeAudioContexts.forEach((ctx) => {
      try {
        if (ctx && ctx.state === 'running') {
          ctx.suspend();
        }
      } catch (err) {
        console.warn('[AudioManager] Erro ao suspender contexto:', err);
      }
    });
  }

  function resumeAppAudio() {
    if (document.hidden) return;
    if (isAppInForeground) return;
    isAppInForeground = true;

    if (window._virtualKeys) {
      window._virtualKeys.PageHidden = false;
    }

    activeAudioContexts.forEach((ctx) => {
      try {
        if (ctx && ctx.state === 'suspended') {
          ctx.resume();
        }
      } catch (err) {
        console.warn('[AudioManager] Erro ao retomar contexto:', err);
      }
    });
  }

  // Eventos de visibilidade (Troca de aba no navegador)
  document.addEventListener('visibilitychange', () => {
    if (document.hidden) {
      suspendAppAudio();
    } else {
      resumeAppAudio();
    }
  });

  // Ciclo de vida da página móvel (Safari iOS, Chrome Android)
  window.addEventListener('pagehide', suspendAppAudio);
  window.addEventListener('pageshow', () => {
    if (!document.hidden) resumeAppAudio();
  });

  // Perda de foco de janela (Minimizar, alternar aplicativo no sistema)
  window.addEventListener('blur', suspendAppAudio);
  window.addEventListener('focus', () => {
    if (!document.hidden) resumeAppAudio();
  });

  // Exporta métodos para depuração ou controle externo caso necessário
  window.AudioManager = {
    suspend: suspendAppAudio,
    resume: resumeAppAudio,
    getActiveContexts: () => activeAudioContexts
  };
})();

