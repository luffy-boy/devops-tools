import request from '@/utils/request'

export function roleAll(data) {
  return request({
    url: 'auth/role/all',
    method: 'get',
    data
  })
}

export function fetchList(params){
  return request({
    url: 'auth/role/list',
    method: 'get',
    params
  })
}

export function create(data,method){
  return request({
    url: 'auth/role',
    method: method,
    data
  })
}

export function getDetail(params){
  return request({
    url: 'auth/role/detail',
    method: 'get',
    params
  })
}

export function del(params){
  return request({
    url: 'auth/role',
    method: 'delete',
    params
  })
}

export function routeEdit(params){
  return request({
    url: 'auth/role/route_edit',
    method: 'post',
    params
  })
}