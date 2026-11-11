export interface components {
  schemas: {
    Ejemplo: {
      id: string;
      nombre: string;
      descripcion: string;
      creadoEn: string;
    };
    ListadoEjemplos: {
      ejemplos: components["schemas"]["Ejemplo"][];
    };
    CrearEjemplo: {
      nombre: string;
      descripcion?: string;
    };
    Problema: {
      type?: string;
      title: string;
      status?: number;
      detail?: string;
    };
  };
}

export interface paths {
  "/api/v1/ejemplos": {
    get: {
      responses: {
        200: {
          content: {
            "application/json": components["schemas"]["ListadoEjemplos"];
          };
        };
        500: {
          content: {
            "application/problem+json": components["schemas"]["Problema"];
          };
        };
      };
    };
    post: {
      requestBody: {
        content: {
          "application/json": components["schemas"]["CrearEjemplo"];
        };
      };
      responses: {
        201: {
          content: {
            "application/json": components["schemas"]["Ejemplo"];
          };
        };
        400: {
          content: {
            "application/problem+json": components["schemas"]["Problema"];
          };
        };
        422: {
          content: {
            "application/problem+json": components["schemas"]["Problema"];
          };
        };
        500: {
          content: {
            "application/problem+json": components["schemas"]["Problema"];
          };
        };
      };
    };
  };
}
