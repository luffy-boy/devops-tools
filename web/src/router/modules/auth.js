/** When your routing table is too long, you can split it into small modules **/

import Layout from '@/layout'

const authRouter = {
  path: '/auth',
  component: Layout,
  redirect: '/auth/admin/list',
  name: 'Auth',
  meta: {
    title: '权限管理',
    icon: 'lock'
  },
  children: [
    {
      path: 'admin',
      component: () => import('@/views/auth/index'),
      name: 'Admin',
      redirect: '/auth/admin/list',
      meta: { title: '管理员' },
      children: [
        {
          path: 'list',
          component: () => import('@/views/auth/admin/list'),
          name: 'AdminList',
          meta: { title: '管理员'},
        },
        {
          path: 'detail',
          component: () => import('@/views/auth/admin/detail'),
          name: 'AdminAdd',
          meta: { title: '管理员新增'},
          hidden: true
        },
        {
          path: 'detail/:id(\\d+)',
          component: () => import('@/views/auth/admin/detail'),
          name: 'AdminEdit',
          meta: { title: '管理员修改'},
          hidden: true
        },
      ]
    },
    {
      path: 'role',
      component: () => import('@/views/auth/index'),
      name: 'Role',
      redirect: '/auth/role/list',
      meta: { title: '角色管理' },
      children: [
        {
          path: 'list',
          component: () => import('@/views/auth/role/list'),
          name: 'RoleList',
          meta: { title: '角色列表'},
        },
        {
          path: 'detail',
          component: () => import('@/views/auth/role/detail'),
          name: 'RoleAdd',
          meta: { title: '新增角色'},
          hidden: true
        },
        {
          path: 'detail/:id(\\d+)',
          component: () => import('@/views/auth/role/detail'),
          name: 'RoleEdit',
          meta: { title: '修改角色'},
          hidden: true
        },
        {
          path: 'route_edit/:id(\\d+)',
          component: () => import('@/views/auth/role/route_edit'),
          name: 'RouteEdit',
          meta: { title: '角色权限'},
          hidden: true
        },
      ]
    },
    {
      path: 'route',
      component: () => import('@/views/auth/index'),
      name: 'Route',
      redirect: '/auth/route/list',
      meta: { title: '路由管理' },
      children: [
        {
          path: 'list',
          component: () => import('@/views/auth/route/list'),
          name: 'RouteList',
          meta: { title: '路由列表'},
        },
        {
          path: 'detail',
          component: () => import('@/views/auth/route/detail'),
          name: 'RouteEdit',
          meta: { title: '新增路由'},
          hidden: true
        },
        {
          path: 'detail/:id(\\d+)',
          component: () => import('@/views/auth/route/detail'),
          name: 'RouteEdit',
          meta: { title: '修改路由'},
          hidden: true
        },
      ]
    },
  ]
}
export default authRouter
