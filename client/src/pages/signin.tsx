/* eslint-disable @typescript-eslint/no-unused-vars */
import axios from 'axios';
import { useState } from 'react';
import toast from "react-hot-toast";
import { BACKEND_URL } from '../lib/constants';
import { fetchToken } from '../lib/fetchToken';
import { Spinner } from '../lib/Spinner';

const SignIn = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
  });
  const [submit, setSubmit] = useState<boolean>(false)

  const toggleForm = () => setIsLogin(!isLogin);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData(prev => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async () => {
    try {
      if (isLogin) {
        setSubmit(true)
        const { data: { data: loginData } } = await axios.post(`${BACKEND_URL}/login`, {  
          email: formData.email, password: formData.password
        })

        console.log(loginData.token)

        localStorage.setItem("token", loginData.token)
        localStorage.setItem("username", loginData.user.username)
        localStorage.setItem("id", loginData.user.id)
        setSubmit(false)

        fetchToken()
        toast.success("Logged In")
        window.location.href = "/"
        return
      }

      setSubmit(true)
      const { data: { data: signupData } } = await axios.post(`${BACKEND_URL}/register`, {  
        username: formData.username, email: formData.email, password: formData.password
      })

      localStorage.setItem("token", signupData.token)
      localStorage.setItem("username", formData.username)
      localStorage.setItem("id", signupData.user.id)
      setSubmit(false)

      toast.success("You are registered")

      fetchToken()
      window.location.href = "/"
    } catch (error) {
      toast.error("error during authentication")
      setSubmit(false)
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100 px-4">
      <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
        <h2 className="text-2xl font-bold mb-6 text-center">
          {isLogin ? 'Login' : 'Sign Up'}
        </h2>

        <form className="space-y-4">
          {!isLogin && (
            <div>
              <label className="block text-sm font-medium text-gray-700">Username</label>
              <input
                type="text"
                name="username"
                value={formData.username}
                onChange={handleChange}
                required
                className="mt-1 w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
              />
            </div>
          )}
          <div>
            <label className="block text-sm font-medium text-gray-700">Email</label>
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              required
              className="mt-1 w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Password</label>
            <input
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
              className="mt-1 w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
            />
          </div>
          <button
            type="button"
            className="w-full items-center flex justify-center bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 transition duration-200"
            onClick={handleSubmit}
          >
            {isLogin ? (
              submit ? (
                <>
                  <Spinner />
                  <span className='ml-2'>Logging in...</span>
                </> 
              ) : "Login"
            ) : (
              submit ? (
                <>
                  <Spinner />
                  <span className='ml-2'>Signing in...</span>
                </>
              ) : "Sign Up"
            )}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-600">
          {isLogin ? "Don't have an account?" : 'Already have an account?'}{' '}
          <button
            onClick={toggleForm}
            className="text-blue-500 hover:underline font-medium"
          >
            {isLogin ? 'Sign up' : 'Login'}
          </button>
        </p>
      </div>
    </div>
  );
};

export default SignIn;
