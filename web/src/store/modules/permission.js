import { asyncRoutes, constantRoutes } from '@/router'
import Layout from '@/layout'
import { 
  getMenu as getCacheMenu
} from '@/utils/auth' // get token from cookie

/**
 * Use meta.role to determine if the current user has permission
 * @param roles
 * @param route
 */
function hasPermission(roles, route) {
  if (route.meta && route.meta.roles) {
    return roles.some(role => route.meta.roles.includes(role))
  } else {
    return true
  }
}

export const loadView = (view) => {
  return (resolve) => require([`@/views/${view}`], resolve)
};

/**
 * 后台查询的菜单数据拼装成路由格式的数据
 * @param routes
 */
export function generaMenu(routes, data) {
  data.forEach(item => {
    const menu = {
      path:  item.path,
      component: item.component === 'layout' ? Layout : loadView(item.component),
      hidden: item.hidden,
      redirect: item.redirect,
      children: [],
      name: item.name,
      meta: item.meta
    }

    if (item.children) {
      generaMenu(menu.children, item.children)
    }
    routes.push(menu)
  })
}
/**
 * Filter asynchronous routing tables by recursion
 * @param routes asyncRoutes
 * @param roles
 */
export function filterAsyncRoutes(routes, roles) {
  const res = []

  routes.forEach(route => {
    const tmp = { ...route }
    if (hasPermission(roles, tmp)) {
      if (tmp.children) {
        tmp.children = filterAsyncRoutes(tmp.children, roles)
      }
      res.push(tmp)
    }
  })

  return res
}

const state = {
  routes: [],
  addRoutes: []
}

const mutations = {
  SET_ROUTES: (state, routes) => {
    state.addRoutes = routes
    state.routes = constantRoutes.concat(routes)
  }
}

const actions = {
  generateRoutes({ commit }, roles) {
    return new Promise(resolve => {
      let tempAsyncRoutes = Object.assign([], asyncRoutes)
      let loadMenuData = []
      const mune = getCacheMenu()
      if ( mune ){
        loadMenuData = mune
      }
      generaMenu(tempAsyncRoutes, loadMenuData)
      let accessedRoutes
      if (roles.includes('admin')) {
        // alert(JSON.stringify(asyncRoutes))
        accessedRoutes = tempAsyncRoutes || []
      } else {
        accessedRoutes = filterAsyncRoutes(tempAsyncRoutes, roles)
      }
      commit('SET_ROUTES', accessedRoutes)
      resolve(accessedRoutes)
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
