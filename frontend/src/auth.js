import React from 'react'
import {createUser, getSession, createSession, deleteSession} from "./api";

class UserInfo {
    constructor(emailAddress, token) {
        this.emailAddress = emailAddress
        this.token = token
    }
}

export class AuthService {
    constructor() {
        this.user = null
    }

    isAuthenticated() {
        return this.user != null
    }

    async checkSession() {
        try {
            this.user = await getSession()
        } catch (e) {
            console.log(e)
        }
    }

    async register(emailAddress, password) {
        const response = await createUser(emailAddress, password)
        this.user = new UserInfo(emailAddress, response.token)
    }
    async login(emailAddress, password) {
        const response = await createSession(emailAddress, password)
        this.user = new UserInfo(emailAddress, response.token)
    }
    async logout() {
        await deleteSession()
        this.user = null
    }
}

export const AuthContext = React.createContext(null)
