import React, { useState, useEffect } from 'react'
import ReactDOM from 'react-dom'
import {BrowserRouter as Router, Switch, Route, Link, Redirect} from "react-router-dom"
import {AuthContext, AuthService} from "./auth"
import VehiclesList from "./components/VehiclesList"
import Calendar from "./components/Calendar"
import Login from "./components/Login"
import Logout from "./components/Logout"
import Register from "./components/Register";
import Home from "./components/Home";

const authService = new AuthService()

function ProtectedRoute({ component, path, ...rest }) {
    const isAuthenticated = authService.isAuthenticated()

    return (
        <Route
            {...rest}
            component={component}
            render={({location}) =>
                isAuthenticated
                    ? children
                    : <Redirect to={{pathname: "/login", state: {from: location}}}/>
            } />
    )
}

const App = ({history}) => {
    const [loading, setLoading] = useState(true)
    const [authenticated, setAuthenticated] = useState(null)

    useEffect(() => {
        if (loading) {
            authService.checkSession()
                .then(() => setLoading(false))
                .then(() => setAuthenticated(authService.isAuthenticated()))
        }
    }, [loading, authenticated])

    if (loading) {
        return <div>Please wait ... </div>
    }

    const logout = async () => {
        await authService.logout()
        setAuthenticated(false)
    }

    return (
        <AuthContext.Provider value={authService}>
            <Router>
                <div>
                    <nav>
                        <ul>
                            <li><Link to="/">Home</Link></li>
                            {authenticated ? <>
                                <li><Link to="/vehicles">Vehicles</Link></li>
                                <li><Link to="/calendar">Calendar</Link></li>
                                <li><a href="" onClick={logout} >Log out</a></li>
                            </> : <>
                                <li><Link to="/register">Register</Link></li>
                                <li><Link to="/login">Login</Link></li>
                            </>}
                        </ul>
                    </nav>

                    <Switch>
                        <ProtectedRoute path="/vehicles" component={VehiclesList} />
                        <ProtectedRoute path="/calendar" component={Calendar} />
                        <ProtectedRoute path="/logout" component={Logout} />

                        <Route path="/register" component={Register} />
                        <Route path="/login" component={Login} />
                        <Route component={Home} />
                    </Switch>
                </div>
            </Router>
        </AuthContext.Provider>
    )
}

const app = document.getElementById('app')
ReactDOM.render(<App />, app)
