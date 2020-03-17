package com.dkrichards.gamesoflife;

import java.io.*;
import java.util.*; 

/**
 * GameProperties is a singleton of game configs.
 */
public class GameProperties {
    private static GameProperties instance = null; 
    private static int worldWidth;
    private static int worldHeight;
    private static double initialSpawnTolerance;
    private static int maxSteps;

    /**
     * Load game configs.
     */
    private GameProperties() throws IOException, FileNotFoundException {
        String rootPath = Thread.currentThread().getContextClassLoader().getResource("").getPath();
        String propertiesPath = rootPath + "gol.properties";
 
        Properties golProps = new Properties();
        golProps.load(new FileInputStream(propertiesPath));

        Random rand = new Random();
 
        String worldWidthProperty = golProps.getProperty("worldWidth");
        String worldHeightProperty = golProps.getProperty("worldHeight");
        String initialSpawnToleranceProperty = golProps.getProperty("initialSpawnTolerance");
        String maxStepsProperty = golProps.getProperty("maxSteps");

        this.worldWidth = Integer.parseInt(worldWidthProperty);
        this.worldHeight = Integer.parseInt(worldHeightProperty);
        this.initialSpawnTolerance = Double.parseDouble(initialSpawnToleranceProperty);
        this.maxSteps = Integer.parseInt(maxStepsProperty);
    }

    /**
     * Return the current instance of GameProperties or create a new one.
     */
    public static GameProperties getInstance() throws IOException { 
        if (instance == null) 
            instance = new GameProperties(); 
  
        return instance; 
    }

    /**
     * Get the world width from properties.
     */
    public int getWorldWidth() {
        return this.worldWidth;
    }

    /**
     * Get the world height from properties.
     */
    public int getWorldHeight() {
        return this.worldHeight;
    }

    /**
     * Get the initial spawn tolerance from properties.
     */
    public double getInitialSpawnTolerance() {
        return this.initialSpawnTolerance;
    }

    /**
     * Get the max number of steps from properties.
     */
    public int getMaxSteps() {
        return this.maxSteps;
    }
} 
