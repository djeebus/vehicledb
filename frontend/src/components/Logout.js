import React, { useEffect, useContext } from 'react'
import { AuthContext } from "../auth";

export default function Logout({history}) {
    const authService = useContext(AuthContext)

    useEffect(() => {
        authService.logout().then(() => history.replace('/login'))
    })

    return <div>Logout</div>
}
