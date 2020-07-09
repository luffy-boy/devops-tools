import request from '@/utils/request'

export function fetchList(params){
  return request({
    url: 'devops/task/list',
    method: 'get',
    params
  })
}
export function create(data,method){
  return request({
    url: 'devops/task',
    method: method,
    data
  })
}

export function auditTask(data){
  return request({
    url: 'devops/task/audit',
    method: 'POST',
    data
  })
}

export function executeTask(params){
  return request({
    url: 'devops/task/execute',
    params
  })
}

export function runningTask(data){
  return request({
    url: 'devops/task/running',
    method: 'POST',
    data
  })
}

export function getDetail(params){
  return request({
    url: 'devops/task/detail',
    method: 'get',
    params
  })
}

export function getGroupList(){
  return request({
    url: 'devops/task/group_list',
    method: 'get',
  })
}

export function del(params){
  return request({
    url: 'devops/task',
    method: 'delete',
    params
  })
}

export function getBanList() {
  return request({
    url: 'devops/task/ban_list',
    method : 'get'
  })
}

export function editBanList(params) {
  return request({
    url: 'devops/task/edit_ban',
    method : 'post',

    params
  })
}

export function fetchLogList(params) {
  return request({
    url: 'devops/task/log',
    params
  })
}

export function logDetail(params) {
  return request({
    url: 'devops/task/log_detail',
    params
  })
}


export function getNotifyData(params){
  return request({
    url: 'devops/task/notify_data',
    method: 'get',
    params
  })
}
