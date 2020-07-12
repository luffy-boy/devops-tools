import request from '@/utils/request'

export function fetchList(params){
  return request({
    url: 'auth/route/list',
    method: 'get',
    params
  })
}

export function routeAll(params){
  return request({
    url: 'auth/route/all',
    method: 'get',
    params
  })
}

export function create(data,method){
  return request({
    url: 'auth/route',
    method: method,
    data
  })
}

export function getDetail(params){
  return request({
    url: 'auth/route/detail',
    method: 'get',
    params
  })
}

export function del(params){
  return request({
    url: 'auth/route',
    method: 'delete',
    params
  })
}

export function editSort(data){
  return request({
    url:  'auth/route/edit_sort',
    method: 'post',
    data
  })
}

export function getAuthMenu() {
  return request({
    url: 'auth/route/menu',
    method: 'get'
  })
}