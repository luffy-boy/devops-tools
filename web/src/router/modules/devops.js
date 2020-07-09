/** When your routing table is too long, you can split it into small modules **/

import Layout from '@/layout'

const devopsRouter = {
  path: '/devops',
  component: Layout,
  redirect: '/devops/servers/list',
  name: 'devops',
  meta: {
    title: '运维',
    icon: 'servers'
  },
  children: [
    {
      path: 'servers',
      component: () => import('@/views/devops/index'),
      name: 'servers',
      redirect: '/devops/servers/list',
      meta: { title: '服务器' },
      children: [
        {
          path: 'list',
          component: () => import('@/views/devops/servers/list'),
          name: 'serversList',
          meta: { title: '服务器'},
        },
        {
          path: 'detail',
          component: () => import('@/views/devops/servers/detail'),
          name: 'serversAdd',
          meta: { title: '新增'},
          hidden: true
        },
        {
          path: 'detail/:id(\\d+)',
          component: () => import('@/views/devops/servers/detail'),
          name: 'serversEdit',
          meta: { title: '修改'},
          hidden: true
        },
      ]
    },
    {
      path: 'task',
      component: () => import('@/views/devops/index'),
      name: 'task',
      redirect: '/devops/task/list',
      meta: { title: '定时任务' },
      children: [
        {
          path: 'list',
          component: () => import('@/views/devops/task/list'),
          name: 'taskList',
          meta: { title: '任务列表'},
        },
        {
          path: 'ban',
          component: () => import('@/views/devops/task/ban'),
          name: 'taskBan',
          meta: { title: '禁用命令'},
        },
        {
          path: 'detail',
          component: () => import('@/views/devops/task/detail'),
          name: 'taskAdd',
          meta: { title: '新增'},
          hidden: true
        },
        {
          path: 'detail/:id(\\d+)',
          component: () => import('@/views/devops/task/detail'),
          name: 'taskEdit',
          meta: { title: '修改'},
          hidden: true
        },
        {
          path: 'log/:id(\\d+)',
          component: () => import('@/views/devops/task/log'),
          name: 'taskLog',
          meta: { title: '任务执行日志'},
          hidden: true
        },
        {
          path: 'log_detail/:id',
          component: () => import('@/views/devops/task/log_detail'),
          name: 'taskLogDetail',
          meta: { title: '任务执行信息'},
          hidden: true
        }
      ]
    },
    {
      path: 'notify',
      component: () => import('@/views/devops/index'),
      name: 'notify',
      redirect: '/devops/notify/list',
      meta: { title: '消息模板' },
      children: [
        {
          path: 'list',
          component: () => import('@/views/devops/notify/list'),
          name: 'notifyList',
          meta: { title: '消息模板'},
        },
        {
          path: 'detail',
          component: () => import('@/views/devops/notify/detail'),
          name: 'notifyAdd',
          meta: { title: '新增'},
          hidden: true
        },
        {
          path: 'detail/:id(\\d+)',
          component: () => import('@/views/devops/notify/detail'),
          name: 'notifyEdit',
          meta: { title: '修改'},
          hidden: true
        }
      ]
    },
  ]
}
export default devopsRouter
