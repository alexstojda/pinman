import {useEffect, useState} from 'react';
import {User} from "./generated";
import {Api} from "./index";
import {useNavigate} from "react-router-dom";

export type AuthOptions = {
  requireAuth?: boolean;
  requireRoles?: string[];
}

export type AuthState = {
  user: User | undefined
}

const api = new Api();

export function useAuth(opts: AuthOptions): AuthState {
  const navigate = useNavigate()

  const [user, setUser] = useState<User>();
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>()

  useEffect(() => {
    api.tryTokenRefresh().then(success => {
      if (success) {
        setIsAuthenticated(true)
      } else {
        setIsAuthenticated(false)
        api.clearJwtToken()
      }
    })
  }, [])

  useEffect(() => {
    if (isAuthenticated === undefined) {
      return
    }

    if (isAuthenticated) {
      api.userApi().usersMeGet().then(r => {
        if (r.status === 200) {
          setUser(r.data.user)
        } else {
          console.error("failed to get user data", r.config.data)
        }
      })
      //TODO: Validate opts.requiredRoles
    } else if (opts.requireAuth) {
      navigate("/login")
    }
  }, [isAuthenticated, opts.requireAuth, navigate])


  return {
    user
  };
}