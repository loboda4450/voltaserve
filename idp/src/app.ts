import '@/infra/env'
import bodyParser from 'body-parser'
import cors from 'cors'
import express from 'express'
import logger from 'morgan'
import passport from 'passport'
import { Strategy as JwtStrategy, ExtractJwt } from 'passport-jwt'
import { URL } from 'url'
import accountRouter from '@/account/router'
import { getConfig } from '@/config/config'
import { errorHandler } from '@/infra/error'
import tokenRouter from '@/token/router'
import userRepo from '@/user/repo'
import userRouter from '@/user/router'
import { client as postgres } from './infra/postgres'

const app = express()

app.use(cors())
app.use(logger('dev'))
app.use(express.json({ limit: '3mb' }))
app.use(express.urlencoded({ extended: true }))
app.use(bodyParser.json())

const tokenConfig = getConfig().token
passport.use(
  new JwtStrategy(
    {
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      secretOrKey: tokenConfig.jwtSigningKey,
      issuer: tokenConfig.issuer,
      audience: tokenConfig.audience,
    },
    async (jwt_payload, done) => {
      try {
        const user = await userRepo.findByID(jwt_payload.sub)
        return done(null, user)
      } catch {
        return done(null, false)
      }
    }
  )
)

app.use('/v1/user', userRouter)
app.use('/v1/accounts', accountRouter)
app.use('/v1/token', tokenRouter)

app.use(errorHandler)

const port = new URL(getConfig().idpURL).port

postgres.connect().then(() => {
  app.listen(port, () => {
    console.log(`Listening on port ${port}`)
  })
})