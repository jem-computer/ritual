// ABOUTME: Parses human-readable schedule strings into cron expressions
// ABOUTME: Supports formats like "daily at 8:00 AM", "every Monday at 9:00 AM", etc.

import cronParser from 'cron-parser';

export interface ParsedSchedule {
  cron: string;
  description: string;
}

export function parseSchedule(schedule: string): ParsedSchedule {
  const normalized = schedule.toLowerCase().trim();
  
  // Daily patterns
  const dailyMatch = normalized.match(/daily at (\d{1,2}):(\d{2})\s*(am|pm)?/);
  if (dailyMatch) {
    const [, hourStr, minute, ampm] = dailyMatch;
    let hour = parseInt(hourStr);
    
    if (ampm === 'pm' && hour !== 12) {
      hour += 12;
    } else if (ampm === 'am' && hour === 12) {
      hour = 0;
    }
    
    return {
      cron: `${minute} ${hour} * * *`,
      description: `Daily at ${hour.toString().padStart(2, '0')}:${minute}`,
    };
  }
  
  // Weekly patterns
  const weeklyMatch = normalized.match(/every (\w+) at (\d{1,2}):(\d{2})\s*(am|pm)?/);
  if (weeklyMatch) {
    const [, dayName, hourStr, minute, ampm] = weeklyMatch;
    let hour = parseInt(hourStr);
    
    if (ampm === 'pm' && hour !== 12) {
      hour += 12;
    } else if (ampm === 'am' && hour === 12) {
      hour = 0;
    }
    
    const dayMap: Record<string, number> = {
      sunday: 0,
      monday: 1,
      tuesday: 2,
      wednesday: 3,
      thursday: 4,
      friday: 5,
      saturday: 6,
    };
    
    const dayNum = dayMap[dayName];
    if (dayNum === undefined) {
      throw new Error(`Invalid day name: ${dayName}`);
    }
    
    return {
      cron: `${minute} ${hour} * * ${dayNum}`,
      description: `Every ${dayName} at ${hour.toString().padStart(2, '0')}:${minute}`,
    };
  }
  
  // Hourly patterns
  const hourlyMatch = normalized.match(/every (\d+) hours?/);
  if (hourlyMatch) {
    const [, hours] = hourlyMatch;
    return {
      cron: `0 */${hours} * * *`,
      description: `Every ${hours} hour${hours === '1' ? '' : 's'}`,
    };
  }
  
  // Every X minutes patterns
  const minuteMatch = normalized.match(/every (\d+) minutes?/);
  if (minuteMatch) {
    const [, minutes] = minuteMatch;
    return {
      cron: `*/${minutes} * * * *`,
      description: `Every ${minutes} minute${minutes === '1' ? '' : 's'}`,
    };
  }
  
  // Monthly patterns (first day of month)
  const monthlyMatch = normalized.match(/monthly at (\d{1,2}):(\d{2})\s*(am|pm)?/);
  if (monthlyMatch) {
    const [, hourStr, minute, ampm] = monthlyMatch;
    let hour = parseInt(hourStr);
    
    if (ampm === 'pm' && hour !== 12) {
      hour += 12;
    } else if (ampm === 'am' && hour === 12) {
      hour = 0;
    }
    
    return {
      cron: `${minute} ${hour} 1 * *`,
      description: `Monthly on the 1st at ${hour.toString().padStart(2, '0')}:${minute}`,
    };
  }
  
  // If no pattern matches, throw an error
  throw new Error(`Unable to parse schedule: "${schedule}"`);
}

// Helper to get next run time
export function getNextRunTime(cron: string): Date | null {
  try {
    const interval = cronParser.parseExpression(cron);
    return interval.next().toDate();
  } catch (error) {
    console.error('Failed to parse cron expression:', error);
    return null;
  }
}

// Validate cron expression
export function validateCron(cron: string): boolean {
  try {
    cronParser.parseExpression(cron);
    return true;
  } catch {
    return false;
  }
}