import request from '@/utils/request'

export function cronSpeNextTime(params){
  return request({
    url: 'devops/cron/next_runtime',
    method: 'get',
    params
  })
}