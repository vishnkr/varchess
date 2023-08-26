import { derived, writable } from 'svelte/store'

const TOKEN_KEY = 'token'

const auth = writable({
  token: window.localStorage.getItem(TOKEN_KEY),
  isGuest: false,
})

const me = derived(auth, ($auth, set) => {
  set({
    ...$auth,
    isAuth: !!$auth.token || $auth.isGuest,
  })
})

function setToken(token: string) {
  auth.update(($auth) => {
    window.localStorage.setItem(TOKEN_KEY, token)

    return {
      ...$auth,
      token,
    }
  })
}

function loginAsGuest() {
  auth.update(($auth) => {
    return { ...$auth, isGuest: true }
  })
}

export default {
  subscribe: me.subscribe,
  setToken,
  loginAsGuest,
}