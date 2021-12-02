import org.jboss.aop.util.AOPLock;

public class Main {
    public static void main(String[] args) {
        AOPLock a = new AOPLock();
        System.out.println(a.getClass().getName());
    }
}
