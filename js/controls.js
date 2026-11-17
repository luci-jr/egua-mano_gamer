/**
 * ============================================================================
 * Égua Mano! Gamer - Sistema de Controles Virtuais & Input
 * Clean Architecture: Camada de Entrada e Interação do Usuário
 * ============================================================================
 * Gerencia o estado de teclas virtuais acessadas pelo WebAssembly Go,
 * o D-Pad clássico estilo NES com thumb-glide analógico suave, e os botões de ação.
 */

(function () {
  'use strict';

  // 1. Objeto Global de Teclas Virtuais lido pelo motor Go
  window._virtualKeys = {
    ArrowLeft: false,
    ArrowRight: false,
    Jump: false,
    JustJump: false,
    ArrowDown: false,
    JustDuck: false,
    Enter: false,
    Escape: false,
    ArrowUp: false,
    Attack: false,
    PageHidden: false,
    JustUnfocused: false,
    Start: false
  };

  function setVirtual(key, state) {
    if (window._virtualKeys) {
      window._virtualKeys[key] = state;
      if (key === 'Jump' && state) {
        window._virtualKeys['JustJump'] = true;
      }
      if (key === 'ArrowDown' && state) {
        window._virtualKeys['JustDuck'] = true;
      }
    }
  }

  // 2. Trava de Orientação Paisagem & Tela Cheia
  function lockLandscape() {
    try {
      if (screen.orientation && screen.orientation.lock) {
        screen.orientation.lock('landscape').catch(() => {});
      }
    } catch (e) {}
  }

  function enterFullScreen() {
    try {
      const doc = window.document;
      const docEl = doc.documentElement;
      const request = docEl.requestFullscreen ||
                      docEl.webkitRequestFullscreen ||
                      docEl.mozRequestFullScreen ||
                      docEl.msRequestFullscreen;

      if (!doc.fullscreenElement && !doc.webkitFullscreenElement && !doc.mozFullScreenElement && !doc.msFullscreenElement) {
        if (request) {
          request.call(docEl).catch(() => {});
        }
      }
    } catch (e) {}
  }

  window.addEventListener('load', lockLandscape);
  document.addEventListener('touchstart', () => {
    lockLandscape();
    enterFullScreen();
  }, { passive: true });
  document.addEventListener('click', () => {
    lockLandscape();
    enterFullScreen();
    window.focus();
  }, { passive: true });

  // 3. Helper de Despacho de Eventos de Teclado Virtuais
  function triggerKey(code, key, keyCode, type) {
    const event = new KeyboardEvent(type, {
      code: code,
      key: key,
      keyCode: keyCode,
      which: keyCode,
      bubbles: true,
      cancelable: true
    });
    window.dispatchEvent(event);
    document.dispatchEvent(event);
    const canvas = document.querySelector('canvas');
    if (canvas) canvas.dispatchEvent(event);
  }

  function bindTouchButton(elemId, code, key, keyCode, vKey) {
    const btn = document.getElementById(elemId);
    if (!btn) return;

    const onStart = (e) => {
      e.preventDefault();
      btn.classList.add('active');
      if (vKey) setVirtual(vKey, true);
      triggerKey(code, key, keyCode, 'keydown');
      if (navigator.vibrate) navigator.vibrate(10);
    };

    const onEnd = (e) => {
      e.preventDefault();
      btn.classList.remove('active');
      if (vKey) setVirtual(vKey, false);
      triggerKey(code, key, keyCode, 'keyup');
    };

    btn.addEventListener('touchstart', onStart, { passive: false });
    btn.addEventListener('touchend', onEnd, { passive: false });
    btn.addEventListener('touchcancel', onEnd, { passive: false });
    btn.addEventListener('mousedown', onStart);
    btn.addEventListener('mouseup', onEnd);
    btn.addEventListener('mouseleave', onEnd);
  }

  // 4. D-PAD NES COM THUMB-GLIDE (Deslize suave 8 direções)
  const dpadContainer = document.getElementById('nes-dpad');
  const dpadArms = {
    Up: document.getElementById('btn-dpad-up'),
    Down: document.getElementById('btn-dpad-down'),
    Left: document.getElementById('btn-dpad-left'),
    Right: document.getElementById('btn-dpad-right')
  };

  let activeDpadState = { Up: false, Down: false, Left: false, Right: false };

  function updateDpadDirection(x, y) {
    if (!dpadContainer) return;
    const rect = dpadContainer.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    const dx = x - centerX;
    const dy = y - centerY;
    const dist = Math.hypot(dx, dy);

    const newState = { Up: false, Down: false, Left: false, Right: false };

    // Zona de toque direcional com zona morta central (12px)
    if (dist > 12) {
      if (Math.abs(dx) > Math.abs(dy) * 0.45) {
        if (dx > 0) newState.Right = true;
        else newState.Left = true;
      }
      if (Math.abs(dy) > Math.abs(dx) * 0.45) {
        if (dy > 0) newState.Down = true;
        else newState.Up = true;
      }
    }

    applyDpadChanges(newState);
  }

  function applyDpadChanges(newState) {
    const keyMap = {
      Up: { code: 'ArrowUp', key: 'ArrowUp', keyCode: 38, vKey: 'ArrowUp' },
      Down: { code: 'ArrowDown', key: 'ArrowDown', keyCode: 40, vKey: 'ArrowDown' },
      Left: { code: 'ArrowLeft', key: 'ArrowLeft', keyCode: 37, vKey: 'ArrowLeft' },
      Right: { code: 'ArrowRight', key: 'ArrowRight', keyCode: 39, vKey: 'ArrowRight' }
    };

    for (const dir in newState) {
      if (newState[dir] !== activeDpadState[dir]) {
        activeDpadState[dir] = newState[dir];
        const m = keyMap[dir];
        if (newState[dir]) {
          dpadArms[dir]?.classList.add('active');
          setVirtual(m.vKey, true);
          triggerKey(m.code, m.key, m.keyCode, 'keydown');
          if (navigator.vibrate) navigator.vibrate(8);
        } else {
          dpadArms[dir]?.classList.remove('active');
          setVirtual(m.vKey, false);
          triggerKey(m.code, m.key, m.keyCode, 'keyup');
        }
      }
    }
  }

  function resetDpad() {
    applyDpadChanges({ Up: false, Down: false, Left: false, Right: false });
  }

  if (dpadContainer) {
    let dpadTouchId = null;

    dpadContainer.addEventListener('touchstart', (e) => {
      e.preventDefault();
      const touch = e.changedTouches[0];
      dpadTouchId = touch.identifier;
      updateDpadDirection(touch.clientX, touch.clientY);
    }, { passive: false });

    dpadContainer.addEventListener('touchmove', (e) => {
      e.preventDefault();
      for (let i = 0; i < e.touches.length; i++) {
        if (e.touches[i].identifier === dpadTouchId) {
          updateDpadDirection(e.touches[i].clientX, e.touches[i].clientY);
          break;
        }
      }
    }, { passive: false });

    const onTouchEnd = (e) => {
      for (let i = 0; i < e.changedTouches.length; i++) {
        if (e.changedTouches[i].identifier === dpadTouchId) {
          dpadTouchId = null;
          resetDpad();
          break;
        }
      }
    };

    dpadContainer.addEventListener('touchend', onTouchEnd, { passive: false });
    dpadContainer.addEventListener('touchcancel', onTouchEnd, { passive: false });

    // Fallback para mouse em testes no navegador desktop
    let isDpadMouseDown = false;
    dpadContainer.addEventListener('mousedown', (e) => {
      isDpadMouseDown = true;
      updateDpadDirection(e.clientX, e.clientY);
    });
    window.addEventListener('mousemove', (e) => {
      if (isDpadMouseDown) updateDpadDirection(e.clientX, e.clientY);
    });
    window.addEventListener('mouseup', () => {
      if (isDpadMouseDown) {
        isDpadMouseDown = false;
        resetDpad();
      }
    });
  }

  // 5. Associação dos Botões de Ação
  // Botão B: Tiro / Ataque com a Baladeira (Space / Attack)
  bindTouchButton('btn-touch-b', 'Space', ' ', 32, 'Attack');

  // Botão A: Pulo com Física Parabólica / Salto Duplo (ArrowUp / Jump)
  bindTouchButton('btn-touch-a', 'ArrowUp', 'ArrowUp', 38, 'Jump');

  // START: Pausar Jogo e Abrir Menu / Confirmar Seleção
  bindTouchButton('btn-touch-start', 'Enter', 'Enter', 13, 'Start');

  // SELECT: Alterna Opções / Herói / Créditos
  bindTouchButton('btn-touch-select', 'ArrowDown', 'ArrowDown', 40, 'ArrowDown');

  // PAUSA (Topo da Tela)
  bindTouchButton('btn-touch-pause', 'Escape', 'Escape', 27, 'Escape');

  // Exporta para escopo global
  window.GameControls = {
    setVirtual: setVirtual,
    triggerKey: triggerKey
  };
})();

