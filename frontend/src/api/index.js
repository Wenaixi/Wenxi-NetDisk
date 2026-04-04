import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  res => res.data,
  err => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(err.response?.data || err)
  }
)

export const authAPI = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
  me: () => api.get('/auth/me')
}

export const fileAPI = {
  list: (params) => api.get('/files', { params }),
  upload: (formData) => api.post('/files/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }),
  getUploadUrl: (data) => api.post('/files/upload-url', data),
  delete: (id) => api.delete(`/files/${id}`),
  rename: (id, name) => api.put(`/files/${id}/rename`, { name }),
  move: (id, folderId) => api.put(`/files/${id}/move`, { folder_id: folderId }),
  updateDescription: (id, description) => api.put(`/files/${id}/description`, { description }),
  getVersions: (id) => api.get(`/files/${id}/versions`),
  restoreVersion: (id, versionId) => api.post(`/files/${id}/versions/${versionId}/restore`),
  initializeUpload: (data) => api.post('/lanzou/upload/init', data),
  completeUpload: (sessionId, data) => api.post(`/lanzou/upload/complete/${sessionId}`, data),
  download: (id) => api.get(`/files/${id}/download`)
}

export const folderAPI = {
  list: (params) => api.get('/folders', { params }),
  create: (data) => api.post('/folders', data),
  delete: (id) => api.delete(`/folders/${id}`),
  rename: (id, name) => api.put(`/folders/${id}`, { name }),
  updateDescription: (id, description) => api.put(`/folders/${id}`, { description }),
  move: (id, parentId) => api.put(`/folders/${id}/move`, { parent_id: parentId })
}

export const shareAPI = {
  create: (data) => api.post('/shares', data),
  list: () => api.get('/shares'),
  get: (token) => api.get(`/shares/${token}`),
  validate: (token, password) => api.post(`/shares/${token}/validate`, { password }),
  delete: (id) => api.delete(`/shares/${id}`)
}

export default api