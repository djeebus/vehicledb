import React, { useContext, useState } from 'react'
import {AuthContext} from "../auth";

export default function Login() {
    const [emailAddress, setEmailAddress] = useState("")
    const [password, setPassword] = useState("")
    const [error, setError] = useState("")
    const authService = useContext(AuthContext)

    async function login() {
        if (!(emailAddress && password)) {
            return
        }

        try {
            await authService.login(emailAddress, password)
            this.props.history.push('/')
        } catch (e) {
            console.log(e)
            setError(e.code)
        }
    }

    return (
        <form>
            <h3>Login</h3>

            {error ? <p className="error">{error}</p> : null}

            <div className="control">
                <label htmlFor="emailAddress">Email Address: </label>
                <input id="emailAddress" type="textbox"
                       onChange={(e) => setEmailAddress(e.target.value)}/>
            </div>

            <div className="control">
                <label htmlFor="password">Password: </label>
                <input id="password" type="password"
                       onChange={(e => setPassword(e.target.value))}/>
            </div>
            <br/>
            <button type="button" onClick={() => login()}>Submit</button>
        </form>
    )
}
