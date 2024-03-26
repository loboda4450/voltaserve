import { combineReducers } from 'redux'
import error from './error'
import files from './files'
import nav from './nav'
import organizations from './organizations'
import uploadsDrawer from './uploads-drawer'

export default combineReducers({
  uploadsDrawer,
  files,
  nav,
  error,
  organizations,
})
