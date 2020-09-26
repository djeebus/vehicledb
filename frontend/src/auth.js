import React from 'react'
import {createUser} from "./api";

class UserInfo {
    constructor(emailAddress, token) {
        self.emailAddress = emailAddress
        self.token = token
    }
}

export class AuthService {
    constructor() {
        self.user = null
    }

    isAuthenticated() {
        return self.user != null
    }

    async register(emailAddress, password) {
        const response = await createUser(emailAddress, password)
        self.user = new UserInfo(emailAddress, response.token)
    }
    async login(emailAddress, password) {
        const response = await createSession(emailAddress, password)
        self.user = new UserInfo(emailAddress, response.token)
    }
    async logout() {
        const response = await deleteSession()
        self.user = null
    }
}

export const AuthContext = React.createContext(null)
