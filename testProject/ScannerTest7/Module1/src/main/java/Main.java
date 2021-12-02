public class Main {
    public static void main(String[] args) {
        MyLogger logger = new MyLogger();
        System.out.println(logger.getClass().getName());
        System.out.println(logger.getClass().getInterfaces().length);
    }
}
