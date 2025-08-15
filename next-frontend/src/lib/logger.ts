const isDevelopment = process.env.NODE_ENV === 'development';
const logLevel = process.env.NEXT_PUBLIC_LOG_LEVEL || 'info';

const shouldLog = (level: string) => {
  const levels = ['debug', 'info', 'warn', 'error'];
  return levels.indexOf(level) >= levels.indexOf(logLevel);
};

export const logger = {
  debug: (message: string, data?: any) => {
    if (isDevelopment && shouldLog('debug')) {
      console.log(`🔍 DEBUG: ${message}`, data || '');
    }
  },
  
  info: (message: string, data?: any) => {
    if (shouldLog('info')) {
      console.log(`ℹ️ INFO: ${message}`, data || '');
    }
  },
  
  warn: (message: string, data?: any) => {
    if (shouldLog('warn')) {
      console.warn(`⚠️ WARN: ${message}`, data || '');
    }
  },
  
  error: (message: string, error?: any) => {
    if (shouldLog('error')) {
      console.error(`❌ ERROR: ${message}`, error || '');
    }
  },

  // Logs específicos para diferentes contextos
  auth: (message: string, data?: any) => {
    logger.info(`🔐 AUTH: ${message}`, data);
  },

  api: (message: string, data?: any) => {
    logger.info(`🌐 API: ${message}`, data);
  },

  route: (message: string, data?: any) => {
    logger.info(`🛣️ ROUTE: ${message}`, data);
  },

  state: (message: string, data?: any) => {
    logger.debug(`📊 STATE: ${message}`, data);
  }
};
