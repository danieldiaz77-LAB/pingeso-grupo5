import { BrowserRouter, Route, Routes } from "react-router";

import { Inicio } from "./features/inicio/Inicio";

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Inicio />} />
        <Route path="*" element={<Inicio />} />
      </Routes>
    </BrowserRouter>
  );
}
