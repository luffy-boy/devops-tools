import request from '@/utils/request'

export function getNotifyData(params){
  return request({
    url: 'devops/notify/conf',
    method: 'get',
    params
  })
}

export function fetchList(params){
  return request({
    url: 'devops/notify/list',
    method: 'get',
    params
  })
}
export function create(data,method){
  return request({
    url: 'devops/notify',
    method: method,
    data
  })
}

export function getDetail(params){
  return request({
    url: 'devops/notify/detail',
    method: 'get',
    params
  })
}

export function del(params){
  return request({
    url: 'devops/notify',
    method: 'delete',
    params
  })
}