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

function ProtectedRoute({ children, ...rest }) {
    return (
        <Route
            {...rest}
            render={({location}) =>
                authService.isAuthenticated()
                    ? children
                    : <Redirect to={{pathname: "/login", state: {from: location}}}/>
            } />
    )
}

const authService = new AuthService()

const App = () => {
    return (
        <AuthContext.Provider value={authService}>
            <Router>
                <div>
                    <nav>
                        <ul>
                            <li><Link to="/">Home</Link></li>
                            {authService.isAuthenticated() ? <>
                                <li><Link to="/vehicles">Vehicles</Link></li>
                                <li><Link to="/calendar">Calendar</Link></li>
                                <li><Link to="/logout">Log out</Link></li>
                            </> : <>
                                <li><Link to="/register">Register</Link></li>
                                <li><Link to="/login">Login</Link></li>
                            </>}
                        </ul>
                    </nav>

                    <Switch>
                        <ProtectedRoute path="/vehicles">
                            <VehiclesList />
                        </ProtectedRoute>
                        <ProtectedRoute path="/calendar">
                            <Calendar />
                        </ProtectedRoute>
                        <ProtectedRoute path="/logout">
                            <Logout />
                        </ProtectedRoute>

                        <Route path="/register">
                            <Register />
                        </Route>
                        <Route path="/login">
                            <Login />
                        </Route>
                        <Route>
                            <Home />
                        </Route>
                    </Switch>
                </div>
            </Router>
        </AuthContext.Provider>
    )
}

const app = document.getElementById('app')
ReactDOM.render(<App />, app)