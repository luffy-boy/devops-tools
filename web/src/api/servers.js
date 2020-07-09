import request from '@/utils/request'

export function fetchList(params){
  return request({
    url: 'devops/servers/list',
    method: 'get',
    params
  })
}
export function create(data,method){
  return request({
    url: 'devops/servers',
    method: method,
    data
  })
}

export function getDetail(params){
  return request({
    url: 'devops/servers/detail',
    method: 'get',
    params
  })
}

export function getGroupList(){
  return request({
    url: 'devops/servers/group_list',
    method: 'get',
  })
}

export function del(params){
  return request({
    url: 'devops/servers',
    method: 'delete',
    params
  })
}