import java.io.*;
import java.util.*;
import com.dkrichards.gamesoflife.*;

class GameOfLife { 
    /**
     * Clear the terminal.
     */
    public static void clearConsole() {  
        System.out.print("\033[H\033[2J");  
        System.out.flush();  
    } 

    /**
     * Play Conway's Game of Life.
     */
    public static void main(String args[]) throws IOException, FileNotFoundException, InterruptedException {
        GameProperties props = GameProperties.getInstance();
        int maxSteps = props.getMaxSteps();

        World world = new World();
        for (int i = 0; i < maxSteps; i++) {
            GameOfLife.clearConsole();
            world.draw();
            world.step();
            Thread.sleep(500);
        }
    }
}