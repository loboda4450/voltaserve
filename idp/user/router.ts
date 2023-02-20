import fs from 'fs/promises'
import http from 'http'
import os from 'os'
import { URL } from 'url'
import { NextFunction, Router, Response } from 'express'
import { body, validationResult } from 'express-validator'
import multer from 'multer'
import passport from 'passport'
import jwt from 'jsonwebtoken'
import { getConfig } from '@/infra/config'
import { ErrorCode, parseValidationError, newError } from '@/infra/error'
import { PassportRequest } from '@/infra/passport-request'
import {
  deleteUser,
  getUser,
  updateEmail,
  updateFullName,
  updatePicture,
  updatePassword,
  UserDeleteOptions,
  UserUpdateEmailOptions,
  UserUpdateFullNameOptions,
  UserUpdatePasswordOptions,
  deletePicture,
  getByPicture,
} from './service'

const router = Router()

router.get(
  '/',
  passport.authenticate('jwt', { session: false }),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      res.json(await getUser(req.user.id))
    } catch (err) {
      next(err)
    }
  }
)

router.post(
  '/update_full_name',
  passport.authenticate('jwt', { session: false }),
  body('fullName').isString().notEmpty().trim().escape().isLength({ max: 255 }),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      res.json(
        await updateFullName(req.user.id, req.body as UserUpdateFullNameOptions)
      )
    } catch (err) {
      next(err)
    }
  }
)

router.post(
  '/update_email',
  passport.authenticate('jwt', { session: false }),
  body('email').isEmail().isLength({ max: 255 }),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      res.json(
        await updateEmail(req.user.id, req.body as UserUpdateEmailOptions)
      )
    } catch (err) {
      next(err)
    }
  }
)

router.post(
  '/update_password',
  passport.authenticate('jwt', { session: false }),
  body('currentPassword').notEmpty(),
  body('newPassword').isStrongPassword(),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      res.json(
        await updatePassword(req.user.id, req.body as UserUpdatePasswordOptions)
      )
    } catch (err) {
      if (err === 'invalid_password') {
        res.status(400).json({ error: err })
        return
      } else {
        next(err)
      }
    }
  }
)

router.post(
  '/update_picture',
  passport.authenticate('jwt', { session: false }),
  multer({
    dest: os.tmpdir(),
    limits: { fileSize: 3000000, fields: 0, files: 1 },
  }).single('file'),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      const user = await updatePicture(
        req.user.id,
        req.file.path,
        req.file.mimetype
      )
      await fs.rm(req.file.path)
      res.json(user)
    } catch (err) {
      next(err)
    }
  }
)

router.post(
  '/delete_picture',
  passport.authenticate('jwt', { session: false }),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      res.json(await deletePicture(req.user.id))
    } catch (err) {
      next(err)
    }
  }
)

router.delete(
  '/',
  passport.authenticate('jwt', { session: false }),
  body('password').isString().notEmpty(),
  async (req: PassportRequest, res: Response, next: NextFunction) => {
    try {
      const result = validationResult(req)
      if (!result.isEmpty()) {
        throw parseValidationError(result)
      }
      await deleteUser(req.user.id, req.body as UserDeleteOptions)
      res.sendStatus(200)
    } catch (err) {
      if (err === 'invalid_password') {
        res.status(400).json({ error: err })
        return
      } else {
        next(err)
      }
    }
  }
)

export default router
