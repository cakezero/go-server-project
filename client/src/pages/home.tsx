/* eslint-disable @typescript-eslint/no-explicit-any */
import axios from 'axios';
import { useState, useEffect } from 'react';
import { BACKEND_URL } from '../lib/constants';
import toast from "react-hot-toast";
import { Spinner } from '../lib/Spinner';

const Home = () => {
  const arrayValues = [{ equation: "", action: "" }]
  const [activeButton, setActiveButton] = useState<number | null>(null);
  const [first, setFirst] = useState<string>("");
  const [second, setSecond] = useState<string>("");
  const [action, setAction] = useState<string>("");
  const [values, setValues] = useState(arrayValues)
  const [submit, setSubmit] = useState(false)
  const [answer, setAnswer] = useState<string | number>("")

  const symbols = ["➕", "➖", "✖️", "➗"]
  const arithmetics = ["add", "subtract", "multiply", "divide"]

  const handleButtonClick = (index: number) => {
    setActiveButton(index === activeButton ? null : index);
    setAction(arithmetics[index]);
  };

  const token = localStorage.getItem("token")
  const username = localStorage.getItem("username");
  const id = localStorage.getItem("id") ? localStorage.getItem("id") : "";

  useEffect(() => {
    (async () => {
      if (!token) return

      const { data: { data: historyData } } = await axios.get(`${BACKEND_URL}/get-history?id=${id}`, {
        headers: { Authorization: `Bearer ${token}` }
      })

      if (historyData === null) {
        setValues(arrayValues)
        return
      }

      setValues(historyData)
    })()
  })

  const handleComputation = async () => {
    try {
      if (first === "" || second === "" || action === "") {
        toast.error("Provide all actions - arithmetic symbol, first and second numbers")
        return
      }

      setSubmit(true)

      const { data: { data } } = await axios.post(`${BACKEND_URL}/perform-action`, { first, id, second, action });

      toast.success("Calculation performed")
      if (!data.equation) {
        setAnswer(data)
        setSubmit(false)
        return
      }

      console.log(data.action, data.answer)

      if (values[0].action === "" && values[0].equation === "") {
        setValues([data])
        setAnswer(data.answer)
        setSubmit(false)
        return
      }

      setValues(prev => [...prev, data])
      setAnswer(data.answer)
      setSubmit(false)
    } catch (error: any) {
      const axiosError = error.response.data.error;

      if (axiosError) {
        toast.error(axiosError)
      } else {
        toast.error("Something went wrong")
      }

      setSubmit(false)
    }
  }

  const logout = async () => {
    try {
      console.log("token:", token)
      const { data } = await axios.get(`${BACKEND_URL}/logout`, {
        headers: { Authorization: `Bearer ${token}` }
      });

      localStorage.removeItem("token");
      localStorage.removeItem("username");
      localStorage.removeItem("id");

      toast.success(data.message);
      window.location.href = "/"
    } catch (error) {
      console.error("Logout failed:", error);
      localStorage.removeItem("token");
      localStorage.removeItem("username");

      localStorage.removeItem("id");
      window.location.href = "/";

      return
    }
  }

  return (
    <div className="min-h-screen bg-white text-gray-800">
      {/* Navbar */}
      <nav className="bg-gray-100 py-3 px-6 flex justify-between items-center border-b">
        <div className="text-xl font-semibold">GO CAL</div>
        <div className="space-x-4">
          {username ?
            <>
              <h2 className="font-bold text-2xl">{username}</h2>
              <button className="" onClick={logout} type="button">Logout</button>
            </>
          :
            <a href="/auth" className="hover:underline">Sign In</a>
          }
        </div>
      </nav>

      {/* Main Container */}
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold text-center mb-8">GO CALCULATOR</h1>

        {/* Input Section */}
        <form className="flex justify-center items-center gap-8 mb-10 flex-wrap">
          {/* Left Large Box */}
          <input
            type="text"
            placeholder="first number"
            onChange={e => setFirst(e.target.value)}
            className="w-24 h-24 p-2 border border-gray-300 rounded text-center text-sm"
          />

          {/* Middle 2x2 Small Buttons */}
          <div className="grid grid-cols-2 gap-2">
            {[0, 1, 2, 3].map((i) => (
              <button
                key={i}
                type="button"
                onClick={() => handleButtonClick(i)}
                className={`w-10 h-10 border rounded text-sm font-semibold transition ${
                  activeButton === i
                    ? 'bg-blue-600 text-white'
                    : 'bg-white text-gray-700 border-gray-400 hover:bg-blue-100'
                }`}
              >
                {symbols[i]}
              </button>
            ))}
          </div>

          {/* Right Large Box */}
          <input
            type="text"
            placeholder="second number"
            onChange={e => setSecond(e.target.value)}
            className="w-24 h-24 p-2 border border-gray-300 rounded text-center text-sm"
          />

          <p> = </p>
          
          {/* Answer Large Box */}
          <input
            type="text"
            placeholder="Answer"
            readOnly
            value={answer}
            className="w-24 h-24 p-2 border border-gray-300 rounded text-center text-sm"
          />

          {/* Submit Button */}
          <button
            type="button"
            onClick={handleComputation}
            className="bg-blue-500 items-center flex justify-center text-white py-2 px-4 rounded hover:bg-blue-700 mt-4 md:mt-0"
          >
            {!submit ? "Calculate" : (
              <>
                <Spinner />
                <span className='ml-2'>Calculating...</span>
              </>
            )}
          </button>
        </form>

        {/* Table */}
        <div className="overflow-x-auto">
          <table className="table-auto w-full border border-gray-200">
            <thead className="bg-gray-200">
              <tr>
                <th className="p-3 text-left">S/N</th>
                <th className="p-3 text-left">Equation</th>
                <th className="p-3 text-left">Action</th>
              </tr>
            </thead>
            <tbody>
              {values[0].action === "" && values[0].equation === "" ? (
                <tr>
                  <td colSpan={3} className="p-3 text-center">
                    No computations performed
                  </td>
                </tr>
              ) : (
                values.map((value, index) => (
                  <tr key={index} className="bg-gray-100">
                    <td className="p-3 text-left">{index + 1}</td>
                    <td className="p-3 text-left">{value.equation}</td>
                    <td className="p-3 text-left">{value.action}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Home;
