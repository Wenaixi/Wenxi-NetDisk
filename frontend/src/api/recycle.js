import api from './index'

// 回收站 API
export const recycleAPI = {
  // 获取回收站列表
  list: () => api.get('/recycle'),

  // 恢复文件
  restoreFile: (id) => api.post(`/recycle/${id}/restore`),

  // 恢复文件夹
  restoreFolder: (id) => api.post(`/recycle/${id}/restore`),

  // 永久删除
  permanentDelete: (id) => api.delete(`/recycle/${id}`),

  // 清空回收站
  clearAll: () => api.delete('/recycle/clear'),
}

export default recycleAPI
