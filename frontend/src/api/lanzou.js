import api from './index'

// 蓝奏云连接状态
export const lanzouAPI = {
  // 连接蓝奏云账号
  connect: (data) => api.post('/lanzou/connect', data),

  // 获取连接状态
  status: () => api.get('/lanzou/status'),

  // 断开连接
  disconnect: () => api.delete('/lanzou/connect'),

  // 获取文件列表 (从蓝奏云)
  listFiles: (params) => api.get('/lanzou/files', { params }),

  // 获取文件夹列表
  listFolders: (params) => api.get('/lanzou/folders', { params }),

  // 创建文件夹
  createFolder: (data) => api.post('/lanzou/folders', data),

  // 删除文件
  deleteFile: (id) => api.delete(`/lanzou/files/${id}`),

  // 删除文件夹
  deleteFolder: (id) => api.delete(`/lanzou/folders/${id}`),

  // 重命名
  rename: (id, data) => api.put(`/lanzou/rename/${id}`, data),

  // 移动
  move: (id, data) => api.put(`/lanzou/${id}/move`, data),

  // 获取下载链接
  getDownloadUrl: (id) => api.get(`/lanzou/files/${id}/url`),

  // 创建分享链接
  createShare: (data) => api.post('/lanzou/share', data),

  // 设置访问密码 (文件或文件夹)
  setAccess: (data) => api.put('/lanzou/access', data),

  // 获取文件详情 (蓝奏云)
  getFileDetail: (id) => api.get(`/lanzou/files/${id}`),

  // 批量删除
  batchDelete: (data) => api.post('/lanzou/batch/delete', data),
}

export default lanzouAPI
