import request from '@/utils/request'

export function login(data) {
  return request({
    url: 'auth/admin/login',
    method: 'post',
    data
  })
}

export function reLogin(data) {
  return request({
    url: 'auth/admin/re_login',
    method: 'post',
    data
  })
}

export function logout() {
  return request({
    url: 'auth/admin/logout',
    method: 'post'
  })
}

export function getInfo() {
  return request({
    url: 'auth/admin',
    method: 'get'
  })
}

export function fetchList(params){
  return request({
    url: 'auth/admin/list',
    method: 'get',
    params
  })
}

export function create(data,method){
  return request({
    url: 'auth/admin',
    method: method,
    data
  })
}

export function getDetail(params){
  return request({
    url: 'auth/admin/detail',
    method: 'get',
    params
  })
}

export function getIndexData(){
  return request({
    url: 'auth/admin/index_data',
    method: 'get'
  })
}


