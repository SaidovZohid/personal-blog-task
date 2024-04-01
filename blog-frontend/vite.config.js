import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import dotenv from 'dotenv';

export default defineConfig(({ mode }) => {
  // Load the appropriate environment file based on the mode (development or production)
  const env = dotenv.config({
    path: `./.env.${mode}`,
  }).parsed;

  return {
    plugins: [react()],
    define: {
      'process.env.API_BASE_URL': JSON.stringify(env.API_BASE_URL), 
    },
  };
});
