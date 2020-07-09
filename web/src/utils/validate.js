/**
 * Created by PanJiaChen on 16/11/18.
 */

/**
 * @param {string} path
 * @returns {Boolean}
 */
export function isExternal(path) {
  return /^(https?:|mailto:|tel:)/.test(path)
}

/**
 * @param {string} str
 * @returns {Boolean}
 */
export function validUsername(str) {
  if (str.length < 5 || str.length > 32){
    return false
  }
  return true
}

/**
 * @param {string} str
 * @returns {Boolean}
 */
export function validPassword(str) {
  if (str.length < 5 || str.length > 32){
    return false
  }
  return true
}
