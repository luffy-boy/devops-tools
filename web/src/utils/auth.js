import Cookies from 'js-cookie'

const TokenKey = 'jwt-token'
const MenuKey = 'menu-list'

export function getToken() {
  return Cookies.get(TokenKey)
}

export function setToken(token) {
  return Cookies.set(TokenKey, token)
}

export function removeToken() {
  return Cookies.remove(TokenKey)
}

export function getMenu() {
  return JSON.parse(localStorage.getItem(MenuKey))
}

export function setMenu(menu) {
  return localStorage.setItem(MenuKey, JSON.stringify(menu));
}

export function removeMenu() {
  return localStorage.removeItem(MenuKey)
}